# Chat Rabbit


## Requirement


这个ProxyController的主要功能是实现一个代理，把客户端请求的API转发到服务端，这个转发要确保隔离客户端的信息，不向服务端传输任何可能和客户端标识相关的信息，对于这些信息，都替换成这个代理的信息，让服务端认为所有的请求都是代理服务发出来的，它查不到任何客户端的信息

