import { FitAddon } from "@xterm/addon-fit";
import { Terminal as XtermTerminal, IDisposable } from "@xterm/xterm";
import { lib } from "libapps";
export declare class Xterm {
    elem: HTMLElement;
    term: XtermTerminal;
    fitAddon: FitAddon;
    resizeListener: () => void;
    dataDisposable: IDisposable | null;
    resizeDisposable: IDisposable | null;
    decoder: lib.UTF8Decoder;
    message: HTMLElement;
    messageTimeout: number;
    messageTimer: number;
    constructor(elem: HTMLElement);
    info(): {
        columns: number;
        rows: number;
    };
    output(data: string): void;
    showMessage(message: string, timeout: number): void;
    removeMessage(): void;
    setWindowTitle(title: string): void;
    setPreferences(value: Record<string, unknown>): void;
    onInput(callback: (input: string) => void): void;
    onResize(callback: (colmuns: number, rows: number) => void): void;
    deactivate(): void;
    reset(): void;
    close(): void;
}
