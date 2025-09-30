using Yarp.ReverseProxy.Configuration;

namespace gateway.YarpUtils;

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
                        new DestinationConfig { Address = "http://auth-service:2000" }
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
                        new DestinationConfig { Address = "http://bridge-service:2017" }
                    },
                },
            },
            new ClusterConfig
            {
                ClusterId = "peer-cluster",
                Destinations = new Dictionary<string, DestinationConfig>
                {
                    {
                        "destination3",
                        new DestinationConfig { Address = "http://peer-service:2013" }
                    },
                },
            },
            new ClusterConfig
            {
                ClusterId = "stream-cluster",
                Destinations = new Dictionary<string, DestinationConfig>
                {
                    {
                        "destination4",
                        new DestinationConfig { Address = "http://streaming-service:2019" }
                    },
                },
            },
            new ClusterConfig
            {
                ClusterId = "resource-cluster",
                Destinations = new Dictionary<string, DestinationConfig>
                {
                    {
                        "destination4",
                        new DestinationConfig { Address = "http://resource-service:2019" }
                    },
                },
            },
        ];
    }
}
