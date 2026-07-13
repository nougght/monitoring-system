export function getGradientColor(colors: string[], percentage: number): string {
    // Clamp percentage between 0 and 100
    percentage = Math.max(0, Math.min(100, percentage));

    // If there's only one color, return it
    if (colors.length === 1) return colors[0];

    // Find the two colors to interpolate between
    const point = (colors.length - 1) * percentage / 100;
    const i = Math.floor(point);
    const t = point - i;

    const color1 = colors[i];
    const color2 = colors[Math.min(i + 1, colors.length - 1)];

    console.log(color1, color2, percentage, colors.length);
    const rgb = [0, 1, 2].map(j => {
        // rgb color values from hex
        const c1 = parseInt(color1.slice(1 + j * 2, 3 + j * 2), 16);
        const c2 = parseInt(color2.slice(1 + j * 2, 3 + j * 2), 16);
        
        // interpolate between values
        return Math.round(c1 * (1 - t) + c2 * t);
    });

    // convert rgb to hex
    return '#' + rgb.map(c => c.toString(16).padStart(2, '0')).join('');
}