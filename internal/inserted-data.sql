INSERT INTO Users (userName, userEmail, userPassword)
VALUES ('Christen Ghaly', 'ChristenGhaly@gmail.com', '123456789');

INSERT INTO Threads (title, createdDate, userId)
VALUES ('Programming', CURRENT_TIMESTAMP, 1);

INSERT INTO Messages (msgbody, createdDate, userId,  threadId)
VALUES ('Hello, this is the first message', CURRENT_TIMESTAMP, 1, 1);