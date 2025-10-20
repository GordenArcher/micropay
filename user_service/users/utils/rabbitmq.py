import json
import time
import logging
from urllib.parse import urlparse

import pika
from pika.exceptions import AMQPConnectionError, ChannelClosedByBroker

from django.conf import settings

logger = logging.getLogger(__name__)

class EventPublisher:
    """
    Simple, resilient publisher:
     - Declares exchange 'events' topic durable
     - Attempts reconnect with exponential backoff for publish
     - Marks messages persistent
    """

    def __init__(self, url=None, exchange="events", exchange_type="topic"):
        self.url = url or settings.RABBITMQ_URL
        self.exchange = exchange
        self.exchange_type = exchange_type
        self._connect()

    def _connect(self):
        params = pika.URLParameters(self.url)
        self.connection = None
        self.channel = None
        for attempt in range(1, 6):
            try:
                self.connection = pika.BlockingConnection(params)
                self.channel = self.connection.channel()
                self.channel.exchange_declare(exchange=self.exchange, exchange_type=self.exchange_type, durable=True)
                logger.info("Connected to RabbitMQ")
                return
            except AMQPConnectionError as e:
                sleep = min(2 ** attempt, 30)
                logger.warning(f"RabbitMQ connection attempt {attempt} failed: {e}. Retrying in {sleep}s.")
                time.sleep(sleep)
        raise ConnectionError("Could not connect to RabbitMQ after retries")

    def publish(self, routing_key, payload, retry=3):
        if not self.channel or self.channel.is_closed:
            try:
                self._connect()
            except Exception as e:
                logger.exception("Failed to reconnect to RabbitMQ")
                raise

        body = json.dumps(payload)
        properties = pika.BasicProperties(content_type="application/json", delivery_mode=2)

        for attempt in range(1, retry + 1):
            try:
                self.channel.basic_publish(
                    exchange=self.exchange,
                    routing_key=routing_key,
                    body=body,
                    properties=properties,
                    mandatory=False
                )
                logger.info("Published event %s: %s", routing_key, body)
                return True
            except (AMQPConnectionError, ChannelClosedByBroker) as e:
                logger.warning("Publish attempt %d failed: %s", attempt, e)
                try:
                    self._connect()
                except Exception:
                    logger.exception("Reconnect failed during publish retry")
                time.sleep(min(2 ** attempt, 10))
        logger.error("Failed to publish event %s after %d attempts", routing_key, retry)
        return False

    def close(self):
        try:
            if self.channel and not self.channel.is_closed:
                self.channel.close()
            if self.connection and not self.connection.is_closed:
                self.connection.close()
        except Exception:
            logger.exception("Error closing RabbitMQ connection")
