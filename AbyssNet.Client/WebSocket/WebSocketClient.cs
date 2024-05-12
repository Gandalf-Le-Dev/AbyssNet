using System.Net.WebSockets;
using System.Text;
using AbyssNet.Client.Common;
using AbyssNet.Client.Common.Interfaces;

namespace AbyssNet.Client.WebSocket;

public class AbyssNetClient(AbyssNetClientWebSocketOptions? options = null) : IAbyssNetClient
{
    private ClientWebSocket _client = null!;

    public async Task InitializeClientAsync()
    {
        Uri uri = new(options?.ServerURI ?? throw new InvalidOperationException("Server URI is not set."));
        var cTs = new CancellationTokenSource();
        cTs.CancelAfter(options.ConnectionTimeout); // TODO: Add to config object (AbyssNetClientOptions)

        _client = new()
        {
            Options =
            {
                KeepAliveInterval = options?.KeepAliveInterval ?? TimeSpan.FromSeconds(10), // Sends pong frames to the server every X seconds
            }
        };

        try
        {
            await _client.ConnectAsync(uri, cTs.Token);
        }
        catch (WebSocketException wsEx)
        {
            Console.WriteLine($"Failed to connect: {wsEx.Message}");
        }
        catch (Exception ex)
        {
            Console.WriteLine($"An error occurred: {ex.Message}");
        }
    }

    public async Task PollAsync()
    {
        if (_client == null)
        {
            throw new InvalidOperationException($"Client is not connected. Consider calling {nameof(InitializeClientAsync)}() first.");
        }

        while (_client.State == WebSocketState.Open)
        {
            var buffer = new byte[1024];
            var result = await _client.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);

            if (result.MessageType == WebSocketMessageType.Close)
            {
                Console.WriteLine("Received close message.");
                await _client.CloseAsync(WebSocketCloseStatus.NormalClosure, "Close message received.", CancellationToken.None);
                break;
            }

            var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
            Console.WriteLine($"Received: {message}");
        }
    }

    public async Task SendAsync(string message)
    {
        if (_client == null)
        {
            throw new InvalidOperationException($"Client is not connected. Consider calling {nameof(InitializeClientAsync)}() first.");
        }

        var buffer = Encoding.UTF8.GetBytes(message);
        await _client.SendAsync(new ArraySegment<byte>(buffer), WebSocketMessageType.Binary, true, CancellationToken.None);
    }
}