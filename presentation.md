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
and how I try to approach writing code now. This comes from some code that I have been working on recently at Weave. It 
will go over downloading a file from a remote server, saving it locally as a temp file and then uploading it to a local 
system. I’ll be using this to help describe the different pieces of clean architecture.

### Entities

The core types that define our business logic. For the example program, it will be the document. All the different 
fields that are required to make the document have meaning will be included there. Things like, the remote URL, the 
filename, any metadata about it, etc. This entity can have methods on it that will help us know more about it, things 
like, give me a file type, give me the name for the file, etc.

### Use Cases

These chain together how we intend to use our entities. The example shows downloading a document from a remote server, 
saving as a temp file, then uploading to the local system. These are the steps to execute for the business logic.

### API’s / Jobs

Entry points into the system. This is how we will begin the process of executing the use cases. In the example program, 
we have a REsT endpoint that will begin the download process.

### DB / Network Devices

These are the external pieces that are required to pull data into the system to allow it to gather all information 
needed to execute the use cases appropriately. For us, we have an external request to download the file, a dependency 
on the filesystem to use a temp file, then a call to the external system to store it.

### Configuration

This is how you wire everything together. I take the approach of DI to link everything together, set up the db with the 
required connection string, get the http client ready to call the external url, etc.

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