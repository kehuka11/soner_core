export interface Detail {
  Line: number
  Message: string
  Cost: number
}

export interface FunctionResult {
  name: string
  score: number
  details: Detail[]
  startLine: number
  endLine: number
  source: string
}

export interface FileResult {
  path: string
  functions: FunctionResult[]
  complexity: number
}

export interface AnalyzeResponse {
  files: FileResult[]
}
