from __future__ import unicode_literals

from django.db import models

# Create your models here.

class DataSave(models.Model):
    lid = models.CharField(max_length=200)
    logname = models.CharField(max_length=200)
    game = models.CharField(max_length=200)
    thread = models.CharField(max_length=200)
    loglevel = models.CharField(max_length=200)
    logtime = models.CharField(max_length=200)
    loginfo = models.TextField()

