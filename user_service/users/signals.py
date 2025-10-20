import json
from django.db.models.signals import post_save
from django.dispatch import receiver
from django.contrib.auth.models import User
from django.db import transaction
from .utils.rabbitmq import EventPublisher

publisher = EventPublisher()  # singleton-like usage

@receiver(post_save, sender=User)
def user_created_handler(sender, instance, created, **kwargs):
    if not created:
        return

    event = {
        "user_id": instance.id,
        "username": instance.username,
        "email": instance.email,
    }

    # ensure we publish only after DB transaction commit
    transaction.on_commit(lambda: publisher.publish("user.created", event))
