#!/usr/bin/env bash
set -e

# load .env if present
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# run migrations (safe idempotent)
python manage.py migrate --noinput

# collect static (WhiteNoise)
python manage.py collectstatic --noinput

# start gunicorn
exec gunicorn user_service.wsgi:application \
  --bind 0.0.0.0:${PORT:-8000} \
  --workers 3 \
  --config gunicorn.conf.py
