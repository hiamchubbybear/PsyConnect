import { Component, OnDestroy, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'chat-page',
  templateUrl: './chatpage.html',
  styleUrl: './chatpage.scss',
  standalone: true,
  imports: [FormsModule],
})
export class ChatComponent implements OnInit, OnDestroy {
  ws!: WebSocket;
  messages: string[] = [];
  inputText = '';

  ngOnInit() {
    const token = localStorage.getItem('access_token');
    const conversationID = 'abc123';

    this.ws = new WebSocket(
      `ws://localhost:8080/ws?conversation_id=${conversationID}`
    );

    this.ws.onopen = () => {
      console.log('Connected');
    };

    this.ws.onmessage = (event) => {
      this.messages.push(event.data);
    };

    this.ws.onclose = () => {
      console.log('Closed');
    };
  }

  ngOnDestroy() {
    this.ws?.close();
  }

  sendMessage() {
    this.ws.send(
      JSON.stringify({
        conversation_id: 'abc123',
        user_id: '123',
        text: this.inputText,
        timestamp: new Date().toISOString(),
      })
    );
    this.inputText = '';
  }
}
