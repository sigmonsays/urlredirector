##
# urlredirector
#
# @file
# @version 0.1

VER = latest
TAG := sigmonsays/urlredirector:$(VER)

help:
	#
	# docker        build docker image
	# dockerpush    push docker image to $(TAG)
	#

docker:
	docker build -t $(TAG) .
dockerpush:
	docker push $(TAG)
# end
