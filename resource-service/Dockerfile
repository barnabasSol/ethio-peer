# ---------- Build stage ----------
FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src
 
COPY ["ResourceService.csproj", "./"]
RUN dotnet restore "./ResourceService.csproj"

# copy the rest and publish
COPY . .
RUN dotnet publish "ResourceService.csproj" -c Release -o /app/publish /p:UseAppHost=false

# ---------- Runtime stage ----------
FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS final
WORKDIR /app

# Serve HTTP on 8080 in the container
ENV ASPNETCORE_URLS=http://+:8080
EXPOSE 8080

COPY --from=build /app/publish .
ENTRYPOINT ["dotnet", "ResourceService.dll"]
