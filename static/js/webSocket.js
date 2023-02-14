const ws = new WebSocket("ws://180.76.194.18:80/ws");
let historyMessage = [];
let historyAnswer = "";
document.createElement('li').style.whiteSpace = 'pre-wrap';
function submitQuestion(){
    let question = document.getElementById('question');
    let loading = document.getElementById('loading');
    let chatBox = document.getElementById('messageList');
    document.getElementById('submit').disabled = true;
    let message = question.value.trim();
    if (message.length === 0) {
        return;
    }

    historyMessage.push(message);
    if (historyMessage.length > 5) {
        // 保留最后五条消息
        historyMessage = historyMessage.slice(-5);
    }
    //
    ws.send(JSON.stringify({history: historyMessage,historyAnswer: historyAnswer}));
    question.value = '';
    let li = document.createElement('li');
    li.style.marginTop = '10px';
    li.style.marginBottom = '10px';
    li.style.textAlign = 'right';
    li.style.color = 'blue';
    let questDiv = document.createElement('div');
    questDiv.innerText = "🤔";
    questDiv.style.textAlign = 'right';
    // 添加文本内容到 li 元素中
    var text = document.createTextNode(message);
    li.appendChild(text);
    // 将 li 元素添加到 ul 元素中
    chatBox.appendChild(questDiv);
    questDiv.appendChild(li);
    chatBox.scrollTop = chatBox.scrollHeight;
    loading.style.display = 'block';
    // 启用提交按钮
    document.getElementById('submit').disabled = false;
}
// form.addEventListener('submit', (event) => {
//
// });

ws.onmessage = (event) => {
    let loading = document.getElementById('loading');
    let chatBox = document.getElementById('messageList');
    loading.style.display = 'none';
    let message = JSON.parse(event.data);
    console.log(message);
    let li = document.createElement('li')
    li.style.marginTop = '10px';
    li.style.marginBottom = '10px';
    li.style.textAlign = 'left';
    li.style.color = 'green';
    let questDiv = document.createElement('div');
    questDiv.innerText = "🤖";
    questDiv.style.textAlign = 'left';
    li.style.whiteSpace = 'pre-wrap';
    // 添加文本内容到 li 元素中
    var text = document.createTextNode(message.message);
    li.innerText = message.message;
    console.log("message"+message.message)
    historyAnswer = message.message;
    chatBox.appendChild(questDiv);
    questDiv.appendChild(li);
    chatBox.scrollTop = chatBox.scrollHeight;
    // 将 li 元素添加到 ul 元素中
};
function clearHistory(){
    historyMessage = [];
    historyAnswer = "";
    let chatBox = document.getElementById('messageList');
    chatBox.innerHTML = "";
}
ws.onerror = function(event) {
    console.error("WebSocket error observed:", event);
};