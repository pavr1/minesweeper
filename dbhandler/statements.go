package dbhandler

const INSERT_USER = "INSERT INTO [dbo].[User] ([UserName],[UserLastName],[Password],[CreatedDate]) VALUES (@p1, @p2, @p3, GetDate())"
const VALIDATE_LOGIN = "SELECT [UserName],[UserLastName],[Password],[CreatedDate] FROM [dbo].[User] WHERE [UserName] = @p1 AND [Password] = @p2"
const CREATE_GAME = "INSERT INTO [dbo].[Game] ([UserId],[CreatedDate],[TimeConsumed],[Status]) VALUES (@p1, GetDate(), @p3, @p4)"
