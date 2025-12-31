<script lang="ts">
  import { onMount } from 'svelte';
  import { Terminal } from '@xterm/xterm';
  import { FitAddon } from '@xterm/addon-fit';
  import '@xterm/xterm/css/xterm.css';

  export let wsUrl: string = 'ws://localhost:8080/ws';

  let terminalElement: HTMLDivElement;
  let ws: WebSocket | null = null;
  let term: Terminal | null = null;

  onMount(() => {
    term = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: '#1a1a1a',
      },
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);

    term.open(terminalElement);
    fitAddon.fit();
    term.writeln('Connecting to server...');

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      term?.writeln('Connected to server!');
      term?.writeln('');
    };

    ws.onmessage = (event) => {
      term?.write(event.data);
    };

    ws.onerror = (error) => {
      term?.writeln('\r\nWebSocket error occurred');
      console.error('WebSocket error:', error);
    };

    ws.onclose = () => {
      term?.writeln('\r\nDisconnected from server');
    };

    term.onData((data) => {
      if (ws?.readyState === WebSocket.OPEN) {
        ws.send(data);
      }
    });

    // const handleResize = () => fitAddon.fit();
    // window.addEventListener('resize', handleResize);

    return () => {
      // window.removeEventListener('resize', handleResize);
      ws?.close();
      term?.dispose();
    };
  });
</script>

<div bind:this={terminalElement} class="border border-base-300 rounded-lg p-2 bg-[#1a1a1a]"></div>
