CREATE PROCEDURE CreateGame (
    @userId INT,
    @timeConsumed INT,
	@id INT OUTPUT
) AS
BEGIN
    INSERT INTO [dbo].[Game] ([UserId],[CreatedDate],[TimeConsumed],[Status]) 
	VALUES 
	(@userId, GetDate(), @timeConsumed, 'Pending')

    SET @id=SCOPE_IDENTITY()
    RETURN  @id
END;