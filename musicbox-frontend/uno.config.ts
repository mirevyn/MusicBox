import {
  defineConfig,
  presetIcons,
  presetWind3,
  transformerDirectives,
  // presetAttributify,
} from "unocss";

export default defineConfig({
  presets: [
    presetWind3(),
    // presetAttributify(),
    presetIcons(),
  ],
  transformers: [transformerDirectives()],
  shortcuts: {
    // 使用 SVG 噪点
    "bg-noise":
      'bg-[url(\'data:image/svg+xml,%3Csvg_viewBox="0_0_200_200"_xmlns="http://www.w3.org/2000/svg"%3E%3Cfilter_id="n"%3E%3CfeTurbulence_type="fractalNoise"_baseFrequency="0.65"_numOctaves="3"/%3E%3C/filter%3E%3Crect_width="100%25"_height="100%25"_filter="url(%23n)"/%3E%3C/svg%3E\')]',
  },
  theme: {
    colors: {
      primary: "#335bf1", // MusicBox 主题蓝
    },
  },
});
