-- Active: 1737475325994@@176.108.252.12@5432@crafa
DROP SCHEMA IF EXISTS data CASCADE;
CREATE SCHEMA data;

CREATE TABLE IF NOT EXISTS data.Projects
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY 
, Name TEXT
, PRIMARY KEY (ID)
) WITH (OIDS=FALSE, FILLFACTOR=90);

CREATE TABLE IF NOT EXISTS data.UserProjects
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY 
, UserID     	INT
, ProjectID     INT 
, PRIMARY KEY (ID)
, FOREIGN KEY (UserID) REFERENCES cmn.Users(ID) DEFERRABLE
, FOREIGN KEY (ProjectID) REFERENCES data.Projects(ID) DEFERRABLE
) WITH (OIDS=FALSE, FILLFACTOR=90);


CREATE TABLE IF NOT EXISTS data.Epics
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY 
, UserID     	INT
, StatusID      INT
, AreaID        INT
, CategoryID    INT
, ProjectID     INT
, Name  TEXT
, Text TEXT
, PRIMARY KEY (ID)
, FOREIGN KEY (UserID) REFERENCES cmn.Users(ID) DEFERRABLE
, FOREIGN KEY (StatusID) REFERENCES ref.Statuses(ID) DEFERRABLE
, FOREIGN KEY (AreaID) REFERENCES ref.Areas(ID) DEFERRABLE
, FOREIGN KEY (CategoryID) REFERENCES ref.Categories(ID) DEFERRABLE
, FOREIGN KEY (ProjectID) REFERENCES data.Projects(ID) DEFERRABLE
) WITH (OIDS=FALSE, FILLFACTOR=90);

CREATE TABLE IF NOT EXISTS data.Requirements
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY 
, UserID     	INT
, ProjectID     INT 
, StatusID      INT
, AreaID        INT
, CategoryID    INT
, IterationID   INT
, EpicID        INT
, Text TEXT
, PRIMARY KEY (ID)
, FOREIGN KEY (UserID) REFERENCES cmn.Users(ID) DEFERRABLE
, FOREIGN KEY (StatusID) REFERENCES ref.Statuses(ID) DEFERRABLE
, FOREIGN KEY (AreaID) REFERENCES ref.Areas(ID) DEFERRABLE
, FOREIGN KEY (CategoryID) REFERENCES ref.Categories(ID) DEFERRABLE
, FOREIGN KEY (IterationID) REFERENCES ref.Iterations(ID) DEFERRABLE
, FOREIGN KEY (EpicID) REFERENCES data.Epics(ID) DEFERRABLE
, FOREIGN KEY (ProjectID) REFERENCES data.Projects(ID) DEFERRABLE
) WITH (OIDS=FALSE, FILLFACTOR=90); 

CREATE TABLE IF NOT EXISTS data.Issues
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY 
, UserID     	INT
, ProjectID     INT 
, StatusID      INT
, AreaID        INT
, CategoryID    INT
, IterationID   INT
, EpicID        INT
, Text TEXT
, PRIMARY KEY (ID)
, FOREIGN KEY (UserID) REFERENCES cmn.Users(ID) DEFERRABLE
, FOREIGN KEY (StatusID) REFERENCES ref.Statuses(ID) DEFERRABLE
, FOREIGN KEY (AreaID) REFERENCES ref.Areas(ID) DEFERRABLE
, FOREIGN KEY (CategoryID) REFERENCES ref.Categories(ID) DEFERRABLE
, FOREIGN KEY (IterationID) REFERENCES ref.Iterations(ID) DEFERRABLE
, FOREIGN KEY (EpicID) REFERENCES data.Epics(ID) DEFERRABLE
, FOREIGN KEY (ProjectID) REFERENCES data.Projects(ID) DEFERRABLE
) WITH (OIDS=FALSE, FILLFACTOR=90); 

CREATE TABLE IF NOT EXISTS data.Conversations
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY
, ProjectID     INT  
, IssueID    	INT
, UserID        INT
, AddTime       TIMESTAMP
, UpdateTime    TIMESTAMP
, Text TEXT
, PRIMARY KEY (ID)
, FOREIGN KEY (UserID) REFERENCES cmn.Users(ID) DEFERRABLE
, FOREIGN KEY (IssueID) REFERENCES data.Issues(ID) DEFERRABLE
, FOREIGN KEY (ProjectID) REFERENCES data.Projects(ID) DEFERRABLE
) WITH (OIDS=FALSE, FILLFACTOR=90); 


CREATE TABLE IF NOT EXISTS data.Features
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY 
, UserID     	INT
, ProjectID     INT 
, StatusID      INT
, AreaID        INT
, CategoryID    INT
, IterationID   INT
, IssueID       INT
, ReqID         INT
, Text TEXT
, PRIMARY KEY (ID)
, FOREIGN KEY (UserID) REFERENCES cmn.Users(ID) DEFERRABLE
, FOREIGN KEY (StatusID) REFERENCES ref.Statuses(ID) DEFERRABLE
, FOREIGN KEY (AreaID) REFERENCES ref.Areas(ID) DEFERRABLE
, FOREIGN KEY (CategoryID) REFERENCES ref.Categories(ID) DEFERRABLE
, FOREIGN KEY (IterationID) REFERENCES ref.Iterations(ID) DEFERRABLE
, FOREIGN KEY (ReqID) REFERENCES data.Requirements(ID) DEFERRABLE
, FOREIGN KEY (IssueID) REFERENCES data.Issues(ID) DEFERRABLE
, FOREIGN KEY (ProjectID) REFERENCES data.Projects(ID) DEFERRABLE
) WITH (OIDS=FALSE, FILLFACTOR=90); 

CREATE TABLE IF NOT EXISTS data.Tasks
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY 
, UserID     	INT
, ProjectID     INT 
, StatusID      INT
, AreaID        INT
, CategoryID    INT
, IterationID   INT
, FeatureID       INT
, Text TEXT
, PRIMARY KEY (ID)
, FOREIGN KEY (UserID) REFERENCES cmn.Users(ID) DEFERRABLE
, FOREIGN KEY (StatusID) REFERENCES ref.Statuses(ID) DEFERRABLE
, FOREIGN KEY (AreaID) REFERENCES ref.Areas(ID) DEFERRABLE
, FOREIGN KEY (CategoryID) REFERENCES ref.Categories(ID) DEFERRABLE
, FOREIGN KEY (IterationID) REFERENCES ref.Iterations(ID) DEFERRABLE
, FOREIGN KEY (FeatureID) REFERENCES data.Features(ID) DEFERRABLE
, FOREIGN KEY (ProjectID) REFERENCES data.Projects(ID) DEFERRABLE
) WITH (OIDS=FALSE, FILLFACTOR=90);

CREATE TABLE IF NOT EXISTS data.Works
(
ID INT 			NOT NULL GENERATED BY DEFAULT AS IDENTITY 
, UserID     	INT
, ProjectID     INT 
, TaskID        INT
, Hours         INT
, Time          TIMESTAMPTZ
, PRIMARY KEY (ID)
, FOREIGN KEY (UserID) REFERENCES cmn.Users(ID) DEFERRABLE
, FOREIGN KEY (TaskID) REFERENCES data.Tasks(ID) DEFERRABLE
, FOREIGN KEY (ProjectID) REFERENCES data.Projects(ID) DEFERRABLE
) WITH (OIDS=FALSE, FILLFACTOR=90);