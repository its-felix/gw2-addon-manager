export enum DateFormat {
  SYSTEM = 'system',
  ISO_8601 = 'iso',
}

export enum ColorScheme {
  SYSTEM = 'system',
  LIGHT = 'light',
  DARK = 'dark',
}

export enum UIDensity {
  COMFORTABLE = 'comfortable',
  COMPACT = 'compact',
}

export interface Preferences {
  dateFormat: DateFormat;
  colorScheme: ColorScheme;
  uiDensity: UIDensity;
}

export type EffectiveColorScheme = Exclude<ColorScheme, ColorScheme.SYSTEM>;

export interface EffectivePreferences extends Preferences {
  effectiveColorScheme: EffectiveColorScheme;
}
