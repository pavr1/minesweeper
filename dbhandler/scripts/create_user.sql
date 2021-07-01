USE [minesweeper]
GO
/****** Object:  StoredProcedure [dbo].[CreateGame]    Script Date: 7/1/2021 05:29:40 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[CreateUser] (
    @name VARCHAR(20),
    @lastName VARCHAR(20),
	@password VARCHAR(20),
	@id INT OUTPUT
) AS
BEGIN
    INSERT INTO [dbo].[User] ([UserName],[UserLastName],[Password],[CreatedDate]) 
	VALUES 
	(@name, @lastName, @password, GETDATE())

    SET @id=SCOPE_IDENTITY()
    RETURN  @id
END;