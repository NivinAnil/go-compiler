import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import Editor from '@monaco-editor/react';
import { FaPlay } from 'react-icons/fa';
import { FileJson, Puzzle, Gem, Code2, Coffee, Flame, Sun, Moon } from 'lucide-react';
import { executeCode,getExecution } from './services/execution';
import { FaGolang } from "react-icons/fa6";
import { FaJava, FaSwift, FaRust, FaPython } from "react-icons/fa";
import { SiGnubash } from "react-icons/si";
import { DiRuby } from "react-icons/di";
import { IoLogoJavascript } from "react-icons/io5";

const Spinner = () => (
  <svg className="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
  </svg>
);

const languages = [
  { name: 'JavaScript', icon: <IoLogoJavascript className="h-5 w-5" />, id: '2', fileName: 'main.js', boilerplate: '// Start coding here' },
  { name: 'Python', icon: <FaPython className="h-5 w-5" />, id: '1', fileName: 'code.py', boilerplate: '# Start coding here' },
  ,
  { name: 'Bash', icon: <SiGnubash className="h-5 w-5" />, id: '3', fileName: 'main.sh', boilerplate: '# Start coding here' },
];

const LanguageSelectionToolbar = ({ selectedLanguage, onLanguageSelect }) => (
  <div className="flex flex-wrap justify-center gap-2 p-2 bg-gray-100 dark:bg-gray-800 rounded-lg shadow-sm">
    {languages.map((lang) => (
      <button
        key={lang.name}
        className={`flex items-center justify-center p-2 rounded-md transition-colors duration-200 ${
          selectedLanguage === lang.id
            ? 'bg-blue-500 text-white'
            : 'bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-600'
        }`}
        onClick={() => onLanguageSelect(lang.id, lang.name, lang.fileName)}
      >
        {lang.icon}
        <span className="ml-2 text-sm font-medium">{lang.name}</span>
      </button>
    ))}
  </div>
);

const FileNameHeader = ({ fileName, onExecute, executing }) => (
  <div className="bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200 px-4 py-2 text-sm font-mono border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
    <span>{fileName}</span>
    <button
      onClick={onExecute}
      disabled={executing}
      className={`px-3 py-1 rounded-md flex items-center justify-center transition-colors duration-200 ${
        executing
          ? 'bg-gray-400 cursor-not-allowed'
          : 'bg-green-500 hover:bg-green-600 text-white'
      }`}
    >
      {executing ? <Spinner /> : <FaPlay className="mr-2" />}
      {executing ? "Executing" : "Run"}
    </button>
  </div>
);

const Terminal = ({ output, isReconnecting }) => (
  <div className="flex flex-col h-full border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
    <div className="bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200 px-4 py-2 text-sm font-mono flex items-center">
      <div className="flex space-x-2 mr-4">
        <div className="w-3 h-3 rounded-full bg-red-500"></div>
        <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
        <div className="w-3 h-3 rounded-full bg-green-500"></div>
      </div>
      <span>Output</span>
    </div>
    <div className="flex-grow bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-200 p-4 font-mono text-sm overflow-auto">
      {output ? (
        <pre className="whitespace-pre-wrap">{output}</pre>
      ) : (
        <p className="text-gray-500 dark:text-gray-400">No output received yet...</p>
      )}
    </div>
  </div>
);

const StdInputBox = ({ stdInput, setStdInput }) => (
  <div className="flex flex-col h-full border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
    <div className="bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200 px-4 py-2 text-sm font-mono flex items-center">
      <span>Standard Input</span>
    </div>
    <textarea
      className="flex-grow bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-200 p-4 font-mono text-sm overflow-auto resize-none focus:outline-none focus:ring-2 focus:ring-blue-500"
      value={stdInput}
      onChange={(e) => setStdInput(e.target.value)}
      placeholder="Enter standard input here..."
    />
  </div>
);

const ThemeToggle = ({ theme, setTheme }) => (
  <button
    className="p-2 rounded-md bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200"
    onClick={() => setTheme(theme === "light" ? "dark" : "light")}
    aria-label="Toggle theme"
  >
    {theme === "light" ? <Moon className="h-5 w-5" /> : <Sun className="h-5 w-5" />}
  </button>
);

const App = () => {
  const [editor, setEditor] = useState(null);
  
  const [output, setOutput] = useState('');
  const [isReconnecting, setIsReconnecting] = useState(false);
  const [executing, setExecuting] = useState(false);
  const [selectedLanguage, setSelectedLanguage] = useState('1'); // Default to Python
  const [editorLanguage, setEditorLanguage] = useState('python');
  const [fileName, setFileName] = useState('main.py'); // Default to Python file name
  const [stdInput, setStdInput] = useState('');
  const [theme, setTheme] = useState('light');
  const [boilerplate, setBoilerplate] = useState('');

  const handleEditorDidMount = (editor, monaco) => {
    setEditor(editor);
  };

  const handleExecute = async () => {
    setExecuting(true);
    const request_id = uuidv4();

    if (editor) {
      const code = editor.getValue();
      const encodedCode = btoa(code);
      const encodedStdInput = btoa(stdInput);

      try {
        // Trigger code execution
        const executeResponse = await executeCode({
          code: encodedCode,
          language_id: +selectedLanguage,
          request_id: request_id,
          stdin: encodedStdInput
        });

        // Start polling for execution result

        setTimeout(async () => {
          await pollForExecutionResult(request_id);
        },300)
        
      } catch (error) {
        setOutput(`Error: ${error.message}`);
      } finally {
        setExecuting(false);
      }
    }
  };

  const pollForExecutionResult = async (request_id) => {
    console.log("Starting polling for execution result...");
  
    let resultReceived = false;
    const pollInterval = 2000; // 2 seconds interval
  
    while (!resultReceived) {
      try {
        console.log(`Polling request with ID: ${request_id}`);
        
        const executionResponse = await getExecution({ request_id });
        console.log("Received response:", executionResponse);
  
        if (executionResponse && executionResponse.output) {
          console.log("Execution result received:", executionResponse.output);
          setOutput(executionResponse.output);
          resultReceived = true;
        } else {
          console.log("No output yet, polling again in 2 seconds...");
          // Delay for the next polling request
          await new Promise((resolve) => setTimeout(resolve, pollInterval));
        }
      } catch (error) {
        console.error("Error during polling:", error);
        setOutput(`Polling Error: ${error.message}`);
        resultReceived = true;
      }
    }
  
    console.log("Polling stopped.");
  };
  
  const handleLanguageSelect = (languageId, languageName, newFileName) => {
    setSelectedLanguage(languageId);
    setEditorLanguage(languageName.toLowerCase());
    setFileName(newFileName);
    setBoilerplate(languages.find((lang) => lang.id === languageId).boilerplate);
  };

  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark');
  }, [theme]);

  return (
    <div className={`min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 transition-colors duration-200 ${theme}`}>
      <div className="container mx-auto p-4 flex flex-col h-screen">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-2xl font-bold">{`</>`} Editor</h1>
          <ThemeToggle theme={theme} setTheme={setTheme} />
        </div>
        <LanguageSelectionToolbar
          selectedLanguage={selectedLanguage}
          onLanguageSelect={handleLanguageSelect}
        />
        <div className="flex flex-col lg:flex-row gap-4 flex-grow mt-4">
          <div className="w-full lg:w-2/3 h-[calc(100vh-16rem)] lg:h-auto">
            <div className="h-full border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden flex flex-col">
              <FileNameHeader 
                fileName={fileName} 
                onExecute={handleExecute}
                executing={executing}
              />
              <div className="flex-grow">
                <Editor
                  language={editorLanguage}
                  defaultValue={boilerplate}
                  onMount={handleEditorDidMount}
                  theme={theme === "light" ? "vs-light" : "vs-dark"}
                  options={{
                    minimap: { enabled: false },
                    fontSize: 14,
                    lineNumbers: 'on',
                    roundedSelection: false,
                    scrollBeyondLastLine: false,
                    readOnly: false,
                  }}
                />
              </div>
            </div>
          </div>
          <div className="w-full lg:w-1/3 h-[calc(100vh-16rem)] lg:h-auto flex flex-col gap-4">
            <div className="h-1/2">
              <Terminal output={output} isReconnecting={isReconnecting} />
            </div>
            <div className="h-1/2">
              <StdInputBox stdInput={stdInput} setStdInput={setStdInput} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default App;