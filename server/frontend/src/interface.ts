export interface RuleProvider {
  behavior: string;
  url: string;
  group: string;
  prepend: boolean;
  name: string;
}

export interface Rule {
  rule: string;
  prepend: boolean;
}

export interface Rename {
  old: string;
  new: string;
}

export interface Config {
  clashType: number;
  subscriptions?: string[];
  proxies?: string[];
  userAgent?: string;
  refresh?: boolean;
  autoTest?: boolean;
  lazy?: boolean;
  nodeList?: boolean;
  ignoreCountryGroup?: boolean;
  useUDP?: boolean;
  template?: string;
  ruleProviders?: RuleProvider[];
  rules?: Rule[];
  sort?: string;
  remove?: string;
  replace?: { [key: string]: string };
}
