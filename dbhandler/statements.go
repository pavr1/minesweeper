package dbhandler

const INSERT_USER = "INSERT INTO [dbo].[User] ([UserName],[UserLastName],[Password],[CreatedDate]) VALUES (@p1, @p2, @p3, GetDate())"
