USE [master]
GO
/****** Object:  Database [minesweeper]    Script Date: 7/4/2021 01:10:03 PM ******/
CREATE DATABASE [minesweeper]
 CONTAINMENT = NONE
 ON  PRIMARY 
( NAME = N'minesweeper', FILENAME = N'C:\Program Files\Microsoft SQL Server\MSSQL15.MSSQLSERVER\MSSQL\DATA\minesweeper.mdf' , SIZE = 8192KB , MAXSIZE = UNLIMITED, FILEGROWTH = 65536KB )
 LOG ON 
( NAME = N'minesweeper_log', FILENAME = N'C:\Program Files\Microsoft SQL Server\MSSQL15.MSSQLSERVER\MSSQL\DATA\minesweeper_log.ldf' , SIZE = 8192KB , MAXSIZE = 2048GB , FILEGROWTH = 65536KB )
 WITH CATALOG_COLLATION = DATABASE_DEFAULT
GO
ALTER DATABASE [minesweeper] SET COMPATIBILITY_LEVEL = 150
GO
IF (1 = FULLTEXTSERVICEPROPERTY('IsFullTextInstalled'))
begin
EXEC [minesweeper].[dbo].[sp_fulltext_database] @action = 'enable'
end
GO
ALTER DATABASE [minesweeper] SET ANSI_NULL_DEFAULT OFF 
GO
ALTER DATABASE [minesweeper] SET ANSI_NULLS OFF 
GO
ALTER DATABASE [minesweeper] SET ANSI_PADDING OFF 
GO
ALTER DATABASE [minesweeper] SET ANSI_WARNINGS OFF 
GO
ALTER DATABASE [minesweeper] SET ARITHABORT OFF 
GO
ALTER DATABASE [minesweeper] SET AUTO_CLOSE OFF 
GO
ALTER DATABASE [minesweeper] SET AUTO_SHRINK OFF 
GO
ALTER DATABASE [minesweeper] SET AUTO_UPDATE_STATISTICS ON 
GO
ALTER DATABASE [minesweeper] SET CURSOR_CLOSE_ON_COMMIT OFF 
GO
ALTER DATABASE [minesweeper] SET CURSOR_DEFAULT  GLOBAL 
GO
ALTER DATABASE [minesweeper] SET CONCAT_NULL_YIELDS_NULL OFF 
GO
ALTER DATABASE [minesweeper] SET NUMERIC_ROUNDABORT OFF 
GO
ALTER DATABASE [minesweeper] SET QUOTED_IDENTIFIER OFF 
GO
ALTER DATABASE [minesweeper] SET RECURSIVE_TRIGGERS OFF 
GO
ALTER DATABASE [minesweeper] SET  DISABLE_BROKER 
GO
ALTER DATABASE [minesweeper] SET AUTO_UPDATE_STATISTICS_ASYNC OFF 
GO
ALTER DATABASE [minesweeper] SET DATE_CORRELATION_OPTIMIZATION OFF 
GO
ALTER DATABASE [minesweeper] SET TRUSTWORTHY OFF 
GO
ALTER DATABASE [minesweeper] SET ALLOW_SNAPSHOT_ISOLATION OFF 
GO
ALTER DATABASE [minesweeper] SET PARAMETERIZATION SIMPLE 
GO
ALTER DATABASE [minesweeper] SET READ_COMMITTED_SNAPSHOT OFF 
GO
ALTER DATABASE [minesweeper] SET HONOR_BROKER_PRIORITY OFF 
GO
ALTER DATABASE [minesweeper] SET RECOVERY FULL 
GO
ALTER DATABASE [minesweeper] SET  MULTI_USER 
GO
ALTER DATABASE [minesweeper] SET PAGE_VERIFY CHECKSUM  
GO
ALTER DATABASE [minesweeper] SET DB_CHAINING OFF 
GO
ALTER DATABASE [minesweeper] SET FILESTREAM( NON_TRANSACTED_ACCESS = OFF ) 
GO
ALTER DATABASE [minesweeper] SET TARGET_RECOVERY_TIME = 60 SECONDS 
GO
ALTER DATABASE [minesweeper] SET DELAYED_DURABILITY = DISABLED 
GO
EXEC sys.sp_db_vardecimal_storage_format N'minesweeper', N'ON'
GO
ALTER DATABASE [minesweeper] SET QUERY_STORE = OFF
GO
USE [minesweeper]
GO
/****** Object:  User [minesweeper]    Script Date: 7/4/2021 01:10:03 PM ******/
CREATE USER [minesweeper] FOR LOGIN [minesweeper] WITH DEFAULT_SCHEMA=[dbo]
GO
/****** Object:  Table [dbo].[Game]    Script Date: 7/4/2021 01:10:03 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Game](
	[GameId] [int] IDENTITY(1,1) NOT NULL,
	[UserId] [int] NOT NULL,
	[CreatedDate] [datetime] NOT NULL,
	[TimeConsumed] [float] NOT NULL,
	[Status] [varchar](20) NOT NULL,
	[Rows] [int] NOT NULL,
	[Columns] [int] NOT NULL,
	[Mines] [int] NOT NULL,
PRIMARY KEY CLUSTERED 
(
	[GameId] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Spot]    Script Date: 7/4/2021 01:10:03 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Spot](
	[SpotId] [int] IDENTITY(1,1) NOT NULL,
	[GameId] [int] NOT NULL,
	[Value] [varchar](20) NULL,
	[X] [int] NOT NULL,
	[Y] [int] NOT NULL,
	[NearSpots] [varchar](max) NOT NULL,
	[Status] [varchar](20) NULL,
PRIMARY KEY CLUSTERED 
(
	[SpotId] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
/****** Object:  Table [dbo].[User]    Script Date: 7/4/2021 01:10:03 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[User](
	[UserId] [int] IDENTITY(1,1) NOT NULL,
	[UserName] [varchar](50) NOT NULL,
	[UserLastName] [varchar](50) NOT NULL,
	[Password] [varchar](50) NOT NULL,
	[CreatedDate] [datetime] NOT NULL,
PRIMARY KEY CLUSTERED 
(
	[UserId] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Game]  WITH CHECK ADD FOREIGN KEY([UserId])
REFERENCES [dbo].[User] ([UserId])
GO
/****** Object:  StoredProcedure [dbo].[CreateGame]    Script Date: 7/4/2021 01:10:03 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[CreateGame] (
    @userId INT,
    @timeConsumed INT,
	@rows INT,
	@columns INT,
	@mines INT
) AS
BEGIN
    INSERT INTO [dbo].[Game] ([UserId],[CreatedDate],[TimeConsumed],[Status],[Rows],[Columns],[Mines]) 
	VALUES 
	(@userId, GetDate(), @timeConsumed, 'Pending', @rows, @columns, @mines)

	DECLARE @id INT
    SET @id=SCOPE_IDENTITY()
    RETURN  @id
END;
GO
/****** Object:  StoredProcedure [dbo].[CreateSpot]    Script Date: 7/4/2021 01:10:03 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[CreateSpot] (
    @gameId INT,
    @value VARCHAR(20),
	@x INT,
	@y INT,
	@nearSpots VARCHAR(50),
	@status VARCHAR(20),
	@id INT OUTPUT
) AS
BEGIN
    INSERT INTO [dbo].[Spot] ([GameId],[Value],[X],[Y],[NearSpots],[Status])
     VALUES (@gameId, @value, @x, @y, @nearSpots, @status)

    SET @id=SCOPE_IDENTITY()
    RETURN  @id
END;
GO
/****** Object:  StoredProcedure [dbo].[CreateUser]    Script Date: 7/4/2021 01:10:03 PM ******/
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
GO
USE [master]
GO
ALTER DATABASE [minesweeper] SET  READ_WRITE 
GO
