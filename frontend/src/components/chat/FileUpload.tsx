'use client';

import { useState, useRef, useCallback } from 'react';

interface FileUploadProps {
    userId: string;
    onUploadComplete?: (files: UploadedFile[]) => void;
    onUploadError?: (error: string) => void;
    maxFiles?: number;
    accept?: string;
}

interface UploadedFile {
    id: string;
    filename: string;
    original_name: string;
    mime_type: string;
    size: number;
    url: string;
}

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001';

export default function FileUpload({
    userId,
    onUploadComplete,
    onUploadError,
    maxFiles = 5,
    accept = 'image/*,video/*,application/pdf,.txt',
}: FileUploadProps) {
    const [isDragging, setIsDragging] = useState(false);
    const [isUploading, setIsUploading] = useState(false);
    const [uploadProgress, setUploadProgress] = useState(0);
    const [previewFiles, setPreviewFiles] = useState<{ file: File; preview: string }[]>([]);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleDragOver = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(true);
    }, []);

    const handleDragLeave = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(false);
    }, []);

    const handleDrop = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(false);
        const files = Array.from(e.dataTransfer.files).slice(0, maxFiles);
        handleFiles(files);
    }, [maxFiles]);

    const handleFileSelect = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            const files = Array.from(e.target.files).slice(0, maxFiles);
            handleFiles(files);
        }
    }, [maxFiles]);

    const handleFiles = async (files: File[]) => {
        // Create previews
        const previews = files.map((file) => ({
            file,
            preview: file.type.startsWith('image/') ? URL.createObjectURL(file) : '',
        }));
        setPreviewFiles(previews);
    };

    const uploadFiles = async () => {
        if (previewFiles.length === 0) return;

        setIsUploading(true);
        setUploadProgress(0);

        try {
            const uploadedFiles: UploadedFile[] = [];

            for (let i = 0; i < previewFiles.length; i++) {
                const { file } = previewFiles[i];
                const formData = new FormData();
                formData.append('file', file);

                const response = await fetch(`${API_BASE}/api/upload?userId=${userId}`, {
                    method: 'POST',
                    body: formData,
                });

                if (!response.ok) {
                    throw new Error(`Failed to upload ${file.name}`);
                }

                const result = await response.json();
                uploadedFiles.push(result);
                setUploadProgress(((i + 1) / previewFiles.length) * 100);
            }

            onUploadComplete?.(uploadedFiles);
            setPreviewFiles([]);
        } catch (error) {
            onUploadError?.(error instanceof Error ? error.message : 'Upload failed');
        } finally {
            setIsUploading(false);
            setUploadProgress(0);
        }
    };

    const removePreview = (index: number) => {
        setPreviewFiles((prev) => prev.filter((_, i) => i !== index));
    };

    const formatFileSize = (bytes: number) => {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
        return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
    };

    return (
        <div className="file-upload">
            {/* Drop Zone */}
            <div
                className={`drop-zone ${isDragging ? 'dragging' : ''}`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                onClick={() => fileInputRef.current?.click()}
            >
                <input
                    ref={fileInputRef}
                    type="file"
                    multiple
                    accept={accept}
                    onChange={handleFileSelect}
                    className="hidden"
                />
                <div className="drop-zone-content">
                    <span className="icon">📎</span>
                    <p>ลากไฟล์มาวางที่นี่ หรือคลิกเพื่อเลือก</p>
                    <p className="hint">รองรับ: รูปภาพ, วิดีโอ, PDF (สูงสุด 10MB)</p>
                </div>
            </div>

            {/* Preview */}
            {previewFiles.length > 0 && (
                <div className="preview-container">
                    {previewFiles.map((item, index) => (
                        <div key={index} className="preview-item">
                            {item.preview ? (
                                <img src={item.preview} alt={item.file.name} className="preview-image" />
                            ) : (
                                <div className="preview-file">
                                    <span>📄</span>
                                </div>
                            )}
                            <div className="preview-info">
                                <span className="filename">{item.file.name}</span>
                                <span className="filesize">{formatFileSize(item.file.size)}</span>
                            </div>
                            <button
                                onClick={(e) => {
                                    e.stopPropagation();
                                    removePreview(index);
                                }}
                                className="remove-btn"
                            >
                                ✕
                            </button>
                        </div>
                    ))}
                    <button
                        onClick={uploadFiles}
                        disabled={isUploading}
                        className="upload-btn"
                    >
                        {isUploading ? `กำลังอัพโหลด ${Math.round(uploadProgress)}%` : 'อัพโหลด'}
                    </button>
                </div>
            )}

            <style jsx>{`
        .file-upload {
          width: 100%;
        }
        .drop-zone {
          border: 2px dashed var(--color-earth-300, #d4c4a8);
          border-radius: 12px;
          padding: 2rem;
          text-align: center;
          cursor: pointer;
          transition: all 0.2s;
          background: var(--color-earth-50, #faf8f5);
        }
        .drop-zone:hover,
        .drop-zone.dragging {
          border-color: var(--color-gold-400, #f5c542);
          background: var(--color-gold-50, #fffef5);
        }
        .drop-zone-content {
          color: var(--color-earth-600, #8b7355);
        }
        .icon {
          font-size: 2.5rem;
          display: block;
          margin-bottom: 0.5rem;
        }
        .hint {
          font-size: 0.8rem;
          opacity: 0.7;
          margin-top: 0.5rem;
        }
        .hidden {
          display: none;
        }
        .preview-container {
          margin-top: 1rem;
          display: flex;
          flex-wrap: wrap;
          gap: 0.5rem;
        }
        .preview-item {
          position: relative;
          width: 100px;
          height: 100px;
          border-radius: 8px;
          overflow: hidden;
          border: 1px solid var(--color-earth-200, #e8dcc8);
        }
        .preview-image {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
        .preview-file {
          width: 100%;
          height: 100%;
          display: flex;
          align-items: center;
          justify-content: center;
          background: var(--color-earth-100, #f5f0e8);
          font-size: 2rem;
        }
        .preview-info {
          position: absolute;
          bottom: 0;
          left: 0;
          right: 0;
          padding: 0.25rem;
          background: rgba(0, 0, 0, 0.7);
          color: white;
          font-size: 0.65rem;
        }
        .filename {
          display: block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        .filesize {
          opacity: 0.8;
        }
        .remove-btn {
          position: absolute;
          top: 4px;
          right: 4px;
          width: 20px;
          height: 20px;
          border-radius: 50%;
          border: none;
          background: rgba(0, 0, 0, 0.5);
          color: white;
          cursor: pointer;
          font-size: 0.7rem;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .remove-btn:hover {
          background: rgba(255, 0, 0, 0.7);
        }
        .upload-btn {
          width: 100%;
          margin-top: 0.5rem;
          padding: 0.75rem 1.5rem;
          background: linear-gradient(135deg, var(--color-gold-400, #f5c542), var(--color-gold-500, #e5b32a));
          color: var(--color-earth-900, #2d2416);
          border: none;
          border-radius: 8px;
          font-weight: 600;
          cursor: pointer;
          transition: transform 0.2s;
        }
        .upload-btn:hover:not(:disabled) {
          transform: translateY(-1px);
        }
        .upload-btn:disabled {
          opacity: 0.7;
          cursor: not-allowed;
        }
      `}</style>
        </div>
    );
}
