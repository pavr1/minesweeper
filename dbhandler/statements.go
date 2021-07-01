package dbhandler

const INSERT_USER = "INSERT INTO [dbo].[User] ([UserName],[UserLastName],[Password],[CreatedDate]) VALUES (@p1, @p2, @p3, GetDate()) SELECT SCOPE_IDENTITY()"
const VALIDATE_LOGIN = "SELECT [UserId],[UserName],[UserLastName],[Password],[CreatedDate] FROM [dbo].[User] WHERE [UserName] = @p1 AND [Password] = @p2"
const CREATE_GAME = `
DECLARE	@return_value int,
		@id int

EXEC	@return_value = [dbo].[CreateGame]
		@userId = @p1,
		@timeConsumed = @p2,
		@id = @id OUTPUT

SELECT	@id as N'@id'
`
