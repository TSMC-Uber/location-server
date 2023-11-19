document.addEventListener('DOMContentLoaded', function () {
    var userID = '123';
    var tripID = '456';

    var socket;
    var connectButton = document.getElementById('connect');
    var disconnectButton = document.getElementById('disconnect');
    var sendMessageButton = document.getElementById('send');
    var messagesTextArea = document.getElementById('messages');
    var messageInput = document.getElementById('message');

    connectButton.addEventListener('click', function () {
        // 替换为你的 WebSocket 服务地址
        socket = new WebSocket('ws://localhost:8080/ws/driver?user_id=' + userID + '&trip_id=' + tripID);

        socket.onopen = function () {
            messagesTextArea.value += 'WebSocket 已连接\n';
            connectButton.disabled = true;
            disconnectButton.disabled = false;
        };

        socket.onmessage = function (event) {
            messagesTextArea.value += '收到消息: ' + event.data + '\n';
        };

        socket.onclose = function () {
            messagesTextArea.value += 'WebSocket 已断开\n';
            connectButton.disabled = false;
            disconnectButton.disabled = true;
        };

        socket.onerror = function (error) {
            messagesTextArea.value += '发生错误: ' + error.message + '\n';
        };
    });

    disconnectButton.addEventListener('click', function () {
        if (socket) {
            socket.close();
        }
    });

    sendMessageButton.addEventListener('click', function () {
        if (socket && socket.readyState === WebSocket.OPEN) {
            var message = messageInput.value;
            socket.send(message);
            messageInput.value = '';
        }
    });
});

// {"user_id": "123","trip_id": "456","location": "{'latitude': 123.456,'longitude': 456.789}"}