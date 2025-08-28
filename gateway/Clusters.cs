using Yarp.ReverseProxy.Configuration;

namespace gateway;

public static class Clusters
{
    public static IReadOnlyList<ClusterConfig> GetClusters()
    {
        return
        [
            new ClusterConfig
            {
                ClusterId = "auth-cluster",
                Destinations = new Dictionary<string, DestinationConfig>
                {
                    {
                        "destination1",
                        new DestinationConfig { Address = "http://localhost:2000" }
                    },
                },
            },
            new ClusterConfig
            {
                ClusterId = "bridge-cluster",
                Destinations = new Dictionary<string, DestinationConfig>
                {
                    {
                        "destination2",
                        new DestinationConfig { Address = "http://localhost:2001" }
                    },
                },
            },
        ];
    }
}
