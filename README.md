# Schwarz IT Code Review Repository

This repository is meant to be used in the onboarding process of Go developers
applying to Schwarz IT.

Code smells and security issues are made on purpose. All the code to be reviewed is
in the [review](review) directory.

## My Comment

> Since the task was not very clearly formulated, I decided not to ask many questions and did it as I feel. 
> I refactored the code into a state that suits me more or less. 
> The most logical thing would have been to get a test-task with a list of tests that should continue to run after my changes,
> but in the current version the tests did not work from the very beginning. 
> Therefore, I decided that I am free to make those changes that seemed logical to me.
> I tried not to come up with the business logic myself, since I don’t know the requirements, 
> so I left almost the logic that was originally, except the `MinBasketValue` field, which was asking for validation.
> <br><br>Thank you!<br><br>
> P.S. I understand that it is better to catch panics, if they happen. For this purpose, I wrote 
> library for myself, for running runners (workers, servers ...), but I decided not to embed this library into this code.
> Here is my code https://github.com/tarmalonchik/server-core
