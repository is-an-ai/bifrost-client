import {useState, useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {CheckAndStartLogin} from "../wailsjs/go/main/App.js";

declare global {
    interface Window {
        go: {
            main: {
                App: {
                    CheckAndStartLogin: () => Promise<void>;
                };
            };
        };
    }
}

function App() {
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const waitForWails = () => {
            return new Promise<void>((resolve) => {
                const check = () => {
                    if (window.go && window.go.main && window.go.main.App) {
                        resolve();
                    } else {
                        setTimeout(check, 100);
                    }
                };
                check();
            });
        };

        const checkAuth = async () => {
            try {
                setIsLoading(true);
                setError(null);
                await waitForWails();
                await CheckAndStartLogin();
            } catch (err) {
                console.error('Auth error:', err);
                setError(err instanceof Error ? err.message : 'Authentication failed');
            } finally {
                setIsLoading(false);
            }
        };

        checkAuth();
    }, []);

    if (isLoading) {
        return (
            <div className="loading">
                <p>Checking authentication...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="error">
                <p>Error: {error}</p>
                <button onClick={() => window.location.reload()}>Retry</button>
            </div>
        );
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">Welcome to Bifrost Client!</div>
        </div>
    );
}

export default App;
