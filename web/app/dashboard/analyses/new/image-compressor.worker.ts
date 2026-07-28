self.onmessage = async (event: MessageEvent) => {
  try {
    const { file, options } = event.data;
    const { maxSizeMB, maxWidthOrHeight, initialQuality } = options;

    // Create image blob URL
    const img = new Image();
    img.src = URL.createObjectURL(file);

    // Wait for image to load
    await new Promise((resolve, reject) => {
      img.onload = resolve;
      img.onerror = reject;
    });

    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');

    if (!ctx) {
      throw new Error('Could not get canvas context');
    }

    // Calculate dimensions maintaining aspect ratio
    let width = img.width;
    let height = img.height;

    if (width > height && width > maxWidthOrHeight) {
      height = Math.round((height * maxWidthOrHeight) / width);
      width = maxWidthOrHeight;
    } else if (height > width && height > maxWidthOrHeight) {
      width = Math.round((width * maxWidthOrHeight) / height);
      height = maxWidthOrHeight;
    }

    canvas.width = width;
    canvas.height = height;

    // Draw image on canvas
    ctx.drawImage(img, 0, 0, width, height);

    // Convert to blob with quality adjustment
    const blob = await new Promise<Blob | null>((resolve, reject) => {
      canvas.toBlob(
        (blob) => {
          resolve(blob);
        },
        file.type || 'image/jpeg',
        initialQuality
      );
    });

    if (!blob) {
      throw new Error('Failed to create blob from canvas');
    }

    // Create File object from blob
    const compressedFile = new File([blob], file.name, {
      type: blob.type,
      lastModified: Date.now()
    });

    // Send back the compressed file
    self.postMessage({ compressedFile });
  } catch (error) {
    // Send back error
    self.postMessage({ error: error instanceof Error ? error.message : String(error) });
  }
};