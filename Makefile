CC=gcc
CXX=g++
CXXFLAGS=-Wall -std=c++11
RM=rm -f
CPPFLAGS=-g $(shell root-config --cflags)
LDFLAGS=-g $(shell root-config --ldflags)
LDLIBS=$(shell root-config --libs)

SRCS=src/main.cpp src/CurlTest.cpp
OBJS=$(subst .cpp,.o,$(SRCS))

all: kuby

kuby: $(OBJS)
	$(CXX) $(LDFLAGS) -o kuby $(OBJS) $(LDLIBS)

depend: .depend

.depend: $(SRCS)
	$(RM) ./.depend
	$(CXX) $(CPPFLAGS) -MM $^>>./.depend;

clean:
	$(RM) $(OBJS)

distclean: clean
	$(RM) *~ .depend

include .depend
