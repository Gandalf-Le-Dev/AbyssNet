using AbyssNet.Client.Common;
using AbyssNet.Client.WebSocket;

AbyssNetClientWebSocketOptions webSocketOptions = new()
{
    ServerURI = "ws://localhost:6666/connect",
    KeepAliveInterval = TimeSpan.FromSeconds(10)
};
var client = new AbyssNetClient(webSocketOptions);

Console.WriteLine($"Trying to connect to {webSocketOptions.ServerURI} ...");
await client.InitializeClientAsync();
Console.WriteLine($"Connected!");

await client.PollAsync();