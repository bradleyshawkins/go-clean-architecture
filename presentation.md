# Thoughts on Clean Architecture in Go

## Intro

I’ve been an engineer in one form or another for almost 10 years. During my time, I’ve written bad code. Code that is 
difficult to understand, difficult to maintain, ugly code, and code that has broken things. This has lead me to ask 
myself the question, how do you write quality code? The kind of code that is fun to read, easy to modify, will work as 
expected, and you don’t get lost reading. During the early part of my career, I obsessed over code quality and 
organization. It has lead me to develop certain opinions on structure and organization.

I love improving my craft and want to be the kind of engineer who is always improving, so I am really looking forwards 
to the conversation at the end. I want to know your opinions and experiences, both the positive and negative. I have 
had discussions with other engineers about what they like and don’t like in the past, so I’m excited to get some more 
here.

## Definition

What do I mean when I say “Clean Architecture?” There are a few things that come to mind as I discuss this. This won’t 
be a comprehensive overview of Clean Architecture, and there will certainly be parts from the book and other sources 
that I leave out, but these are the pieces that resonate with me and parts that I try to focus on when writing the code.
I’ll be using a diagram that I quite enjoyed while reading about Clean Architecture

At the end of the presentation I have an example that goes over ways that I have written code and seen code written, 
and how I try to approach writing code now. This comes from some code that I have been working on recently at Weave. 
It will go over downloading a file from a remote server, saving it locally as a temp file and then uploading it to a 
local system. I’ll be using this to help describe the different pieces of clean architecture.

### Entities

The core types that define our business logic. These will vary service to service. Examples User, Account, Email, 
Document, etc. One thing that I try to do is identify my entities when working in a service. This are types that will 
be shared among the rest of the service. Notice how in the diagram the entities are in the center of the circle. This 
indicates that they should have no external dependencies other than basic types. Things like a URL. This is not the 
place to place calls to the database. That will be handled in another layer.

For the example program, it will be the document. All of the different fields that are required to make the document 
have meaning will be included there. Things like, the remote URL, the filename, any metadata about it, etc. This entity 
can have methods on it that will help us know more about it, things like, give me a file type, give me the name for the 
file, etc. These methods will make working with the entities much easier, they will also encompass all of the business 
logic for what the entities can do. Using this pattern makes writing unit tests around business logic much simpler. 
Your business logic is no longer spread throughout incoming requests and outgoing requests to get data.

### Use Cases

These chain together how we intend to use our entities and allow us to create rules for how to use them. This layer is 
one layer out from the entities. This indicates that use cases has an awareness of entities, but not of things external 
to that. Your database technology or communication protocols should not matter here. Obviously you can have a need to 
access a database or a third party service and retrieve new information, if that is the case, I prefer to hide the 
concrete implementation behind an interface. I’ll talk a bit more about how I like to use interfaces during the example. 
For now, suffice it to say, interfaces should be used to hide away the concrete implementation to allow for simpler
unit testing and ease of changing out concrete implementations.

Similarly to entities, creating unit tests for these should be fairly straight forward. The main difference here is 
that you will want to create some sort of mock.

The example shows downloading a document from a remote server, saving as a temp file, then uploading to the local 
system. These are the steps to execute for the business logic.

### API’s / Jobs

Entry points into the system. This is how we will begin the process of executing the use cases. Common implementations 
that I have used are REST, gRPC and a cron job. Each of these allow a use case to be triggered. One preference that I 
have is to have the payload for the request to be separate from the entity. It can be common when first starting to 
have the types contain nearly the same fields. I prefer to create a struct for the request, then a struct for the 
entity. It’s a little bit more code, but it allows you to be able to evolve the types separately.

I tend to think of testing at this layer as integration testing. I like to spin up all of my dependencies up in docker 
containers via docker-compose and test against them. That allows me to verify that all of my database queries are 
working as they should, verify that my migration scripts are correct and running as expected, etc.

In the example program, we have a REsT endpoint that will begin the download process and “upload” process for our document.

### DB / Network Devices

These are the external pieces that are required to pull data into the system to allow it to gather all information 
needed to execute the use cases appropriately. This will often take shape in a database, 3rd party services, either 
internal to your company or external, caches, filesystems, etc. I try to keep my services as agnostic of what 
technology is being used as possible. My service shouldn’t care if I am running Postgres or Mysql as the database, Redis
or Memcached as the cache or communicating via REST or gRPC. The only thing it should care about is the data. Once it 
has the data, it can operate on it as expected.

For us, we have an external request to download the file, a dependency on the filesystem to use a temp file, then a 
call to the external system to store it.

### Configuration

This is how you wire everything together. I take the approach of DI to link everything together. I use main to wire all
of my structs together. I'll create everything there, inject the dependencies and get the system wired up and working
together. You'll see an example of this during the example project.

## Benefits and Drawbacks

Every decision you make in software double-sided, you get some benefits and some drawbacks. Clean Architecture is no 
different. Here are some of the benefits and drawbacks that I’ve identified. This isn’t a comprehensive list, but a few 
that first came to mind for me.

Benefits
* Reusable Code
* Testability
* Can be easier to reason about
* I find myself thinking more about what the overall goal is, how it should be achieved, and how to break things apart
* Readability

Drawbacks
* More file structure
* Can get lost when looking for specific pieces of code
* More code
* More time thinking

## Example

Let’s go through the example now.