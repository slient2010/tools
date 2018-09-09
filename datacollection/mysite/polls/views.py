#from django.shortcuts import render
from django.shortcuts import render
from django.http import HttpResponse
#from polls.libs.mclient import *
from libs.mclient import client

from polls.tasks import addsum

# Create your views here.
def index(request):
    return HttpResponse("Hello, world. You're at the polls index.")

def add(request):
#    a = request.GET['a']
#    b = request.GET['b']
    a = request.GET.get('a', 0)
    b = request.GET.get('b', 0)
    #c = int(a) + int(b)
    c = addsum.delay(int(a),int(b))
    return HttpResponse(str(c.get()))
  

def home(request):
    strings = u"I'm learning Django now"
    return render(request, 'home.html', {'testring': strings})

def forlist(request):
    tlists = ["html", "CSS", "jQuery", "Python", "Django"]
    return render(request, 'test.html', {'testlist': tlists})

def forloop(request):
    List = map(str, range(1, 11))
    return render(request, 'test.html', {'List': List})

def showdata(request):
    testdata = client.getdata()
    if testdata == 0:
        data = '' 
    else:
        data = eval(testdata)[0]['hits']['hits']
#        print data[0]['hits']['hits'][0]['_source']
    
    return render(request, 'data.html', {'data': data})



