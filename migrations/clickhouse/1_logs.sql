CREATE TABLE
    IF NOT EXISTS logs.goods (
        Id int,
        ProjectId int,
        Name VARCHAR(255),
        Description VARCHAR(255),
        Priority int,
        Removed bool,
        EventTime datetime
    ) ENGINE = MergeTree()
ORDER BY
    (Id, ProjectId, Name);
