CREATE PROCEDURE CreateGame (
    @userId INT,
    @timeConsumed INT,
    @rows INT,
    @columns INT,
    @mines INT,
	@id INT OUTPUT
) AS
BEGIN
    INSERT INTO [dbo].[Game] ([UserId],[CreatedDate],[TimeConsumed],[Status],[Rows],[Columns],[Mines]) 
	VALUES 
	(@userId, GetDate(), @timeConsumed, 'Pending', @rows, @columns, @mines)

    SET @id=SCOPE_IDENTITY()
    RETURN  @id
END;