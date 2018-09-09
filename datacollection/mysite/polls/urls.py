from django.shortcuts import render
from django.conf.urls import url

from . import views

urlpatterns = [
    url(r'^$', views.index, name='index'),
    url(r'^home/', views.home, name='home'),
    url(r'^forlist/', views.forlist, name='forlist'),
    url(r'^forloop/', views.forloop, name='forloop'),
    url(r'^showdata/', views.showdata, name='showdata'),
    url(r'^add/', views.add, name='add'),
]
