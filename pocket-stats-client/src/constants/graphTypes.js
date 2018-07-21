const ChartType = {
  PIE: 'PIE',
  LINE: 'LINE'
};

const GraphTypes = {
  ARTICLES_READ: {
    name: "Articles Read",
    type: ChartType.LINE
  },
  ARTICLES_ADDED: {
    name: "Articles Added",
    type: ChartType.LINE
  },
  WORDS_READ: {
    name: "Words Read",
    type: ChartType.LINE
  },
  WORDS_ADDED: {
    name: "Words Added",
    type: ChartType.LINE
  },
  TAGS_READ: {
    name: "Tags Read",
    type: ChartType.PIE
  }
};


export {GraphTypes, ChartType};
