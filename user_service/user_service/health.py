from django.http import JsonResponse
from django.db import connections
from django.db.utils import OperationalError
import os

def healthz(request):
    # simple checks: DB connectivity (quick), env
    db_conn = connections["default"]
    try:
        c = db_conn.cursor()
    except OperationalError:
        return JsonResponse({"status": "fail", "reason": "db_unreachable"}, status=503)
    return JsonResponse({"status": "ok", "service": os.getenv("SERVICE_NAME", "user-service")})
