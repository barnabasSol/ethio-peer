FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /app

COPY ["ResourceService.csproj", "./"]
RUN dotnet restore "./ResourceService.csproj"

COPY . .
RUN dotnet publish "ResourceService.csproj" -c Release -o /app/publish /p:UseAppHost=false


FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app

COPY --from=build /app/publish .
ENTRYPOINT ["dotnet", "ResourceService.dll"]
