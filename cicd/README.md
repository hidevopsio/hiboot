# hicicd

I’m trying not to reinvent the wheel as long as there is an open source CI/CD tool that meets my needs that any developers can build, test, and deploy the web application from the source code without any configurations, and the answer is NO. That’s the hicicd comes to my mind.

Microservices architecture is very popular in modern web application development, build, test, and deploy the application is pretty straight forward, however, to connect, manage, and secure microservices is much harder compare to monolithic web application.

As we transition our applications towards a distributed architecture with microservices deployed across a distributed network, many new challenges await us.

Technologies like containers and container orchestration platforms like Kubernetes or Openshift solve the deployment of our distributed applications quite well, but are still catching up to addressing the service communication necessary to fully take advantage of distributed applications, such as dealing with:

* Unpredictable failure modes
* Verifying end-to-end application correctness
* Unexpected system degradation
* Continuous topology changes
* The use of elastic/ephemeral/transient resources

Today, developers are responsible for taking into account these challenges, and do things like:

* Circuit breaking and Bulkheading (e.g. with Netfix Hystrix)
* Timeouts/retries
* Service discovery (e.g. with Eureka)
* Client-side load balancing (e.g. with Netfix Ribbon)

Another challenge is each runtime and language addresses these with different libraries and frameworks, and in some cases there may be no implementation of a particular library for your chosen language or runtime.

hicicd integrated a new project called [Istio](https://github.com/istio/istio) which will solve many of these challenges and result in a much more robust, reliable, and resilient application in the face of the new world of dynamic distributed applications.

hicicd is a service that can deploy on Kubernetes or Openshift, My mission is to build a tool that should be:

* Easy to use
* Zero configuration
* It’s still configurable for advanced user.

## Git workflow

Below, we outline one of the more common Git workflows that core developers use. Other Git workflows are also valid.

### Fork the main repository

* Go to https://github.com/hidevopsio/hicicd
* Click the "Fork" button (at the top right)

### Clone your fork

The commands below require that you have $GOPATH set ($GOPATH docs). We highly recommend you put Istio's code into your GOPATH. Note: the commands below will not work if there is more than one directory in your $GOPATH.

```bash
export GITHUB_USER=your-github-username
mkdir -p $GOPATH/src/github.com/hidevopsio
cd $GOPATH/src/github.com/hidevopsio
git clone https://github.com/$GITHUB_USER/hicicd
cd hicicd
git remote add upstream 'https://github.com/hidevopsio/hicicd'
git config --global --add http.followRedirects 1
```

### Create a branch and make changes

```bash
git checkout -b my-feature
# Make your code changes
```

### Keeping your fork in sync

```bash
git fetch upstream
git rebase upstream/master
```

Note: If you have write access to the main repositories (e.g. github.com/hidevopsio/hicicd), you should modify your Git configuration so that you can't accidentally push to upstream:

###


