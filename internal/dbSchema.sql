-- Users table
CREATE TABLE Users (
    userId INTEGER NOT NULL PRIMARY KEY,
    userName TEXT NOT NULL,
    userEmail TEXT NOT NULL UNIQUE,
    userPassword TEXT NOT NULL
);

-- Threads table
CREATE TABLE Threads (
    threadId INTEGER NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    createdDate DATETIME NOT NULL,
    userId INTEGER NOT NULL, 
    FOREIGN KEY (userId) REFERENCES Users (userId)
);

-- Messages table
CREATE TABLE Messages (
    msgId INTEGER NOT NULL PRIMARY KEY,
    msgBody TEXT NOT NULL,
    createdDate DATETIME NOT NULL,
    userId INTEGER NOT NULL, 
    threadId INTEGER NOT NULL, 
    FOREIGN KEY (userId) REFERENCES Users (userId),
    FOREIGN KEY (threadId) REFERENCES Threads (threadId)
);

-- Sessions Table
CREATE TABLE Sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry DATETIME NOT NULL
);

CREATE INDEX sessions_expiry_indx ON Sessions(expiry);