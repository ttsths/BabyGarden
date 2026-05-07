import { useState, useCallback } from 'react';

interface UseSpeechRecognitionOptions {
  onResult?: (transcript: string) => void;
  onError?: (error: Error) => void;
  lang?: string;
}

interface SpeechRecognitionEvent {
  resultIndex: number;
  results: Array<Array<{ transcript: string; isFinal: boolean }>>;
}

interface SpeechRecognitionError {
  error: string;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const WindowWithSpeech = window as Record<string, any>;

const getSpeechRecognition = (): new () => {
  start: () => void;
  stop: () => void;
  lang: string;
  continuous: boolean;
  interimResults: boolean;
  onstart: (() => void) | null;
  onresult: ((event: SpeechRecognitionEvent) => void) | null;
  onerror: ((event: SpeechRecognitionError) => void) | null;
  onend: (() => void) | null;
} | undefined => {
  return WindowWithSpeech.SpeechRecognition || WindowWithSpeech.webkitSpeechRecognition;
};

/**
 * 语音识别 Hook (Web Speech API)
 * @param options 配置选项
 */
export function useSpeechRecognition({
  onResult,
  onError,
  lang = 'zh-CN',
}: UseSpeechRecognitionOptions = {}) {
  const [isListening, setIsListening] = useState(false);
  const [transcript, setTranscript] = useState('');
  const [isSupported] = useState(() => !!getSpeechRecognition());

  const startListening = useCallback(() => {
    const SpeechRecognitionClass = getSpeechRecognition();
    
    if (!SpeechRecognitionClass) {
      onError?.(new Error('浏览器不支持语音识别'));
      return;
    }

    const recognition = new SpeechRecognitionClass();
    recognition.lang = lang;
    recognition.continuous = false;
    recognition.interimResults = true;

    recognition.onstart = () => {
      setIsListening(true);
    };

    recognition.onresult = (event: SpeechRecognitionEvent) => {
      let finalTranscript = '';
      let interimTranscript = '';

      for (let i = event.resultIndex; i < event.results.length; i++) {
        const transcript = event.results[i][0].transcript;
        if (event.results[i].isFinal) {
          finalTranscript += transcript;
        } else {
          interimTranscript += transcript;
        }
      }

      const currentTranscript = finalTranscript || interimTranscript;
      setTranscript(currentTranscript);
      onResult?.(currentTranscript);
    };

    recognition.onerror = (event: SpeechRecognitionError) => {
      const error = new Error(event.error || '语音识别错误');
      onError?.(error);
      setIsListening(false);
    };

    recognition.onend = () => {
      setIsListening(false);
    };

    recognition.start();
  }, [lang, onResult, onError]);

  const stopListening = useCallback(() => {
    const SpeechRecognitionClass = getSpeechRecognition();
    
    if (!SpeechRecognitionClass) return;

    const recognition = new SpeechRecognitionClass();
    recognition.stop();
    setIsListening(false);
  }, []);

  const resetTranscript = useCallback(() => {
    setTranscript('');
  }, []);

  return {
    isListening,
    transcript,
    isSupported,
    startListening,
    stopListening,
    resetTranscript,
  };
}
