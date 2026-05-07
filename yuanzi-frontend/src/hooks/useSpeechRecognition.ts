import { useState, useCallback } from 'react';

interface UseSpeechRecognitionOptions {
  onResult?: (transcript: string) => void;
  onError?: (error: Error) => void;
  lang?: string;
}

interface SpeechRecognitionAlternative {
  transcript: string;
  confidence: number;
}

interface SpeechRecognitionResult {
  readonly isFinal: boolean;
  readonly length: number;
  [index: number]: SpeechRecognitionAlternative;
}

interface SpeechRecognitionEvent {
  readonly resultIndex: number;
  readonly results: {
    readonly length: number;
    [index: number]: SpeechRecognitionResult;
  };
}

interface SpeechRecognitionErrorEvent {
  error: string;
}

interface SpeechRecognitionInstance {
  start: () => void;
  stop: () => void;
  lang: string;
  continuous: boolean;
  interimResults: boolean;
  onstart: ((this: SpeechRecognitionInstance) => void) | null;
  onresult: ((this: SpeechRecognitionInstance, event: SpeechRecognitionEvent) => void) | null;
  onerror: ((this: SpeechRecognitionInstance, event: SpeechRecognitionErrorEvent) => void) | null;
  onend: ((this: SpeechRecognitionInstance) => void) | null;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const WindowWithSpeech = window as Record<string, any>;

function getSpeechRecognitionConstructor(): new () => SpeechRecognitionInstance | undefined {
  return WindowWithSpeech.SpeechRecognition || WindowWithSpeech.webkitSpeechRecognition;
}

/**
 * 语音识别 Hook (Web Speech API)
 */
export function useSpeechRecognition({
  onResult,
  onError,
  lang = 'zh-CN',
}: UseSpeechRecognitionOptions = {}) {
  const [isListening, setIsListening] = useState(false);
  const [transcript, setTranscript] = useState('');
  const [isSupported] = useState(() => !!getSpeechRecognitionConstructor());

  const startListening = useCallback(() => {
    const SpeechRecognitionClass = getSpeechRecognitionConstructor();

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
        const result = event.results[i];
        const transcript = result[0].transcript;
        if (result.isFinal) {
          finalTranscript += transcript;
        } else {
          interimTranscript += transcript;
        }
      }

      const currentTranscript = finalTranscript || interimTranscript;
      setTranscript(currentTranscript);
      onResult?.(currentTranscript);
    };

    recognition.onerror = (event: SpeechRecognitionErrorEvent) => {
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
    const SpeechRecognitionClass = getSpeechRecognitionConstructor();

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
