# Technical Questions and Quizzes

This project was made with the proposal to learn some basics of **GO Language** (Golang)  and to implement the application backend based on the idea of microservices.
___
## What does this application do?
The application had services like, ***quiz***, and data. The first one,that is the ***auth***, is responsible for **Signing Up**, **Logging in**, and **Generating Token**. This token is important for being able to use other services like APIs and Quiz services.

The api service was in charge of responding to the user's questions. These questions are multiple choice, and they have categories and difficulty levels. The user can filter the questions returned by using the parameters of the request. They are **Category**, **Difficulty** and  **Number** and are the only the filters can be used in the request.

The quiz service is responsible for returning quizzes using the request of api..They can be personalized using the same filters that were mentioned previously. This service can Generate and Save Quizzes, assign quizzes to users, and Return Quizzes that he creates and others that are assigned to him.

The data service is responsible for getting and saving questions from other public APIs.

---

Now I will explain some decisions I made in the implementation of the idea and why I made them. As well, I will give my opinion on whether they were the best decisions for the project or in real life.

---

## Why? Are they really the best choice?

- Why Gin?
    -  swd

- Why does (almost) every service have a single independent database?
    - swd

- Why SQL and not nSQL?
    - swd

**NOTE:** They are more choices can be explain, although the previous one, in my opinion, are the most relevant 

---

Obviously, this project can be more efficient and readable, and other things should be done in the future.

---
To do?

- Improve the readability of the quiz (go files)

- Make more efficient the auth, api, and (obviously) data and quiz (add timeout for the requests, remove the sleep in the ***data***...)

- Add features like:
    - Respond to the Quiz and save your response.
    - Get the result of responses.
    - Develop/Configure the Orchestrator using Kubernetes(k8s) or other technology 
        - Develop API Gateway
        - Develop Message Broker (mainly because of quiz and api interaction)
    - FrontEnd
    - Generate and Save the results and analytics of the responses to the quizzes for each user.










