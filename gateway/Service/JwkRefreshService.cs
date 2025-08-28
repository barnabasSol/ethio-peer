namespace gateway.Service;

public class JwksRefreshService(IJwksProvider jwksProvider) : BackgroundService
{
    private readonly IJwksProvider _jwksProvider = jwksProvider;
    private readonly TimeSpan _refreshInterval = TimeSpan.FromHours(24);

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        await _jwksProvider.RefreshKeysAsync();

        while (!stoppingToken.IsCancellationRequested)
        {
            await Task.Delay(_refreshInterval, stoppingToken);
            await _jwksProvider.RefreshKeysAsync();
        }
    }
}
