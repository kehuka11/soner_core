import { useLocation, useNavigate } from 'react-router-dom'
import type { FunctionResult, FileResult } from '../types'


function Details() {
    const location = useLocation()
    const navigate = useNavigate()
    const { file, fn, repoUrl } = location.state as { file: FileResult; fn: FunctionResult; repoUrl: string } || {}

    if (!fn) {
        return (
            <div className="card">
                <p>No function details available.</p>
                <button onClick={() => navigate('/')}>Go Back</button>
            </div>
        )
    }

    const getComplexityClass = (score: number) => {
        if (score > 15) return 'complexity-high'
        if (score > 5) return 'complexity-medium'
        return 'complexity-low'
    }

    const getGitHubUrl = () => {
        // Remove .git extension if present
        const cleanUrl = repoUrl.replace(/\.git$/, '')
        // Construct URL: repo/blob/HEAD/path#Lstart-Lend
        return `${cleanUrl}/blob/HEAD/${file.path}#L${fn.startLine}-L${fn.endLine}`
    }

    return (
        <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
            <button onClick={() => navigate('/')} style={{ marginBottom: '1rem' }}>
                &larr; Back to Results
            </button>

            <div className="card">
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem', borderBottom: '1px solid var(--border-color)', paddingBottom: '1rem' }}>
                    <div>
                        <h2 style={{ margin: 0 }}>func {fn.name}</h2>
                        <p style={{ color: '#94a3b8', margin: '0.5rem 0 0 0' }}>{file.path}</p>
                    </div>
                    <div style={{ textAlign: 'right' }}>
                        <div className={getComplexityClass(fn.score)} style={{ fontSize: '1.5em', fontWeight: 'bold' }}>
                            Score: {fn.score}
                        </div>
                        <a
                            href={getGitHubUrl()}
                            target="_blank"
                            rel="noopener noreferrer"
                            style={{ display: 'inline-block', marginTop: '0.5rem', color: '#646cff', textDecoration: 'none' }}
                        >
                            View on GitHub &rarr;
                        </a>
                    </div>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 300px', gap: '2rem' }}>
                    <div>
                        <h3>Source Code</h3>
                        <pre style={{
                            background: '#1e293b',
                            padding: '1rem',
                            borderRadius: '8px',
                            overflowX: 'auto',
                            fontSize: '0.9em',
                            lineHeight: '1.5'
                        }}>
                            <code>{fn.source}</code>
                        </pre>
                    </div>

                    <div>
                        <h3>Complexity Details</h3>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                            {fn.details.map((d, i) => (
                                <div key={i} style={{ background: 'rgba(0,0,0,0.2)', padding: '0.8rem', borderRadius: '4px' }}>
                                    <div style={{ fontWeight: 'bold', color: '#ef4444' }}>+{d.Cost}</div>
                                    <div>{d.Message}</div>
                                    <div style={{ fontSize: '0.8em', color: '#94a3b8' }}>Line {d.Line}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Details
