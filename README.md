# LIA: Part 2

> Due: October 1 \
> Weight: 20%

For this course's Learning Integration Assessment (LIA), you will build
a REST server for an online forum; i.e., a discussion site where people
can hold conversations in the form of posted messages. The assessment is
divided into three parts. The second part focuses on setting up the
database and displaying dynamic data. You will find the requirements
below.

## Requirements

-   Discussion threads should have an ID, a title, and an author (who is
    a user). Messages should have an ID, a body, and an author (who is a
    user). Messages belong to a single discussion thread, but a
    discussion thread can have many messages. Both discussion threads
    and messages should also have a date which corresponds to when they
    were created. Users should have an ID, a username, an email, and a
    password. These informations should be stored in an SQLite database,
    which should be normalized. Operations that deal with the database
    should be defined in a separate `models` package (as seen in class).
    You should have a `Thread` model, a `Message` model, and a `User`
    model, each with the appropriate methods.

-   The homepage should display the 10 latest discussion threads. For
    each thread, the title, the date, the author's username as well as
    the latest message should be shown. The username of the message's
    author should be displayed as well as the first 100 characters of
    its body. Clicking on a discussion thread should bring you to the
    thread's page.

-   Each discussion thread should have its own page where users can read
    all of its messages. The title of the thread, its date of creation
    as well as its author should be shown. Messages should be listed
    from oldest to newest. The username of the message's author, as well
    as its body and date should be displayed.

-   When called, handlers responsible for creating a new discussion
    thread and posting a new message should add the appropriate
    information to the database. Said information can be hard-coded in
    the handlers for now (as done in class). When posting a new message,
    however, the thread to which it belongs should be sourced from the
    URL.

-   Each HTTP request should be logged by the server using a middleware.
    The log entry should include the method as well as the URI of the
    request. Security headers should also be set for each response (see
    the textbook).

-   Runtime errors should be handled correctly. If the server has an
    issue rendering a template, a response with the proper status code
    should be sent to the client.

-   Templates should be cached as seen in class so that the server does
    not have to parse them for every request.

-   The current year should be dynamically displayed in the footer of
    each page.

## Note

Everything needed to produce this assignment can be found in *Let's Go*
chapters 4 (minus section 9), 5 (minus section 6) and 6 (minus sections
4 and 5).

## Submission

The project must be submitted in a repository using GitHub Classroom. To
create the repository, click [here][], and accept the assignment.

[here]: https://classroom.github.com/a/SBwiyKdA

Your first commit should only contain the code you submitted for part 1.
If you use the solution provided by the teacher, your first commit
should only contain the code in the solution. The message of the commit
should be `Add part 1 code`.

## Assessment criteria [20]

-   Readability [5]

    -   code is free of unused variables and functions
    -   use of whitespace/indentation is tidy and consistent
    -   long lines (~80 chars) are split
    -   whitespace is used to visually support logical separation
    -   variable/function names are consistent and descriptive
    -   constants are used instead of hard-coded values
    -   comments are present where warranted, prose is correct and
        helpful
    -   interfaces are well documented
    -   inline comments are used sparingly where needed to decipher
        dense/complex lines

-   Language conventions [2]

    -   no unnecessary use of obscure constructs
    -   standard language features are used appropriately

-   Program design [10]

    -   requirements are met
    -   program flow is decomposed into manageable, logical pieces
    -   function interfaces are clean and well encapsulated
    -   appropriate algorithms are used, and coded cleanly
    -   common code is unified, not duplicated
    -   errors are handled correctly and provide context

-   Data structures [3]

    -   data structures are appropriate
    -   no redundant storage/copying of data
    -   no global variables

