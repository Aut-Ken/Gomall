import { Component, ErrorInfo, ReactNode } from 'react';


interface Props {
    children: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

export default class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
        error: null,
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('Uncaught error:', error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div style={{
                    padding: '40px',
                    textAlign: 'center',
                    fontFamily: 'system-ui, -apple-system, sans-serif'
                }}>
                    <h1 style={{ color: '#ef4444' }}>页面出错了</h1>
                    <p style={{ color: '#666', marginBottom: '20px' }}>
                        我们遇到了一些问题，无法为您显示内容。
                    </p>
                    <details style={{ whiteSpace: 'pre-wrap', textAlign: 'left', background: '#f8fafc', padding: '20px', borderRadius: '8px', overflow: 'auto' }}>
                        {this.state.error && this.state.error.toString()}
                    </details>
                    <button
                        onClick={() => window.location.reload()}
                        style={{
                            marginTop: '20px',
                            padding: '10px 20px',
                            background: '#3b82f6',
                            color: 'white',
                            border: 'none',
                            borderRadius: '6px',
                            cursor: 'pointer'
                        }}
                    >
                        刷新页面
                    </button>
                </div>
            );
        }

        return this.props.children;
    }
}
