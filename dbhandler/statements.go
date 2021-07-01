package dbhandler

const CREATE_USER = `
DECLARE	@return_value int,
		@id int

EXEC	@return_value = [dbo].[CreateUser]
		@name = @p1,
		@lastName = @p2,
		@password = @p3,
		@id = @id OUTPUT

SELECT	@id as N'@id'
`
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
