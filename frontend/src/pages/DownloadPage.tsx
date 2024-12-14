import { useState } from 'react';

export default function DownloadPage() {
    const [code, setCode] = useState('');

    const handleDownload = async () => {
        if (code.length === 5) {
            const response = await fetch(`/api/download/${code}`);
            if (response.ok) {
                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'file';
                document.body.appendChild(a);
                a.click();
                a.remove();
            } else if (response.status === 404) {
                alert('File not found. Please check the code and try again.');
            } else {
                alert('An error occurred. Please try again later.');
            }
        } else {
            alert('Please enter a valid 5-digit code.');
        }
    };

    return (
        <div>
            <h1>Download Page</h1>
            <input
                type="text"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                maxLength={5}
                placeholder="Enter 5-digit code"
                className="border p-2"
            />
            <button onClick={handleDownload} className="ml-2 p-2 bg-blue-500 text-white">
                Download
            </button>
            <br />
            <a href="/upload" className="text-fuchsia-600">To upload page</a>
        </div>
    );
}