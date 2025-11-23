import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import type { AnalyzeResponse, FileResult } from '../types'


function Home() {
    const [repoUrl, setRepoUrl] = useState('')
    const [loading, setLoading] = useState(false)
    const [results, setResults] = useState<FileResult[] | null>(null)
    const [error, setError] = useState<string | null>(null)
    const navigate = useNavigate()

    const handleAnalyze = async () => {
        if (!repoUrl) return
        setLoading(true)
        setError(null)
        setResults(null)

        try {
            const response = await fetch('http://localhost:8081/api/analyze', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ repoUrl }),
            })

            if (!response.ok) {
                throw new Error(`Error: ${response.statusText}`)
            }

            const data: AnalyzeResponse = await response.json()
            setResults(data.files)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An unknown error occurred')
        } finally {
            setLoading(false)
        }
    }

    const getComplexityClass = (score: number) => {
        if (score > 15) return 'complexity-high'
        if (score > 5) return 'complexity-medium'
        return 'complexity-low'
    }

    const handleFunctionClick = (file: FileResult, fn: any) => {
        navigate('/details', { state: { file, fn, repoUrl } })
    }

    return (
        <>
            <h1>SonerCore</h1>
            <p style={{ color: '#94a3b8', marginBottom: '2rem' }}>
                Cognitive Load Analysis for Go Repositories
            </p>

            <div className="card" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'transparent', border: 'none', boxShadow: 'none' }}>
                <input
                    type="text"
                    placeholder="https://github.com/username/repo"
                    value={repoUrl}
                    onChange={(e) => setRepoUrl(e.target.value)}
                    disabled={loading}
                />
                <button onClick={handleAnalyze} disabled={loading || !repoUrl}>
                    {loading ? 'Analyzing...' : 'Analyze Repository'}
                </button>
            </div>

            {error && (
                <div className="card" style={{ borderColor: '#ef4444', color: '#ef4444' }}>
                    {error}
                </div>
            )}

            {results && (
                <div className="file-list">
                    {results.length === 0 ? (
                        <div className="card">No Go files found or analyzed.</div>
                    ) : (
                        results.map((file) => (
                            <div key={file.path} className="card">
                                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem', borderBottom: '1px solid var(--border-color)', paddingBottom: '0.5rem' }}>
                                    <span style={{ fontWeight: 'bold' }}>{file.path}</span>
                                    <span className={getComplexityClass(file.complexity)}>
                                        Total Complexity: {file.complexity}
                                    </span>
                                </div>

                                {file.functions && file.functions.length > 0 ? (
                                    <div style={{ display: 'grid', gap: '0.5rem' }}>
                                        {file.functions.map((fn) => (
                                            <div
                                                key={fn.name}
                                                style={{ background: 'rgba(0,0,0,0.2)', padding: '0.5rem', borderRadius: '4px', cursor: 'pointer' }}
                                                onClick={() => handleFunctionClick(file, fn)}
                                                className="function-item"
                                            >
                                                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                                                    <span>func {fn.name}</span>
                                                    <span className={getComplexityClass(fn.score)}>
                                                        {fn.score}
                                                    </span>
                                                </div>
                                                {fn.details && fn.details.length > 0 && (
                                                    <div style={{ fontSize: '0.8em', color: '#94a3b8', marginTop: '0.5rem', paddingLeft: '1rem' }}>
                                                        {fn.details.length} complexity contributors
                                                    </div>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                ) : (
                                    <div style={{ color: '#64748b', fontStyle: 'italic' }}>No functions analyzed</div>
                                )}
                            </div>
                        ))
                    )}
                </div>
            )}
        </>
    )
}

export default Home
