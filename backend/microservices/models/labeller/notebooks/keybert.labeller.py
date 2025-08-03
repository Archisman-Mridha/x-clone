import marimo

__generated_with = "0.14.10"
app = marimo.App(width="full", app_title="KeyBERT Post Labeller")


@app.cell
def _():
    post = """
    Spent an amazing night in the heart of Copenhagen, with 2 of my lovely colleagues ❤️
    @codingsafari and @johannesholzel.

    We talked primarily about CloudNative Distributed Systems, Data stream processing using Apache
    Flink and building Recommendation Engines.

    And dammnn! Those pastries were real good....😋

    #copenhagen #tech #distributed-systems #big-data #data-stream-processing
    """

    import re

    def preprocess_post(post):
        replacements = [
            # replace newline characters with simple spaces.
            (r'\r\n|\r|\n', ' '),

            # remove mentions and hashtags.
            (r'@[^\s]+', ''),
            (r'#[^\s]+', '')
        ]

        for replacement in replacements:
            old, new = replacement
            post = re.sub(old, new, post)

        return post

    preprocessed_post = preprocess_post(post)
    print(preprocessed_post)
    return (preprocessed_post,)


@app.cell
def _(preprocessed_post):
    import spacy
    import keybert

    # KeyBERT supports quite a few embedding models. Having the option to choose embedding models
    # allow you to leverage pre-trained embeddings that suit your use-case.

    def run_keyBERT_With_spacy():
        spacy.prefer_gpu()
        spacy_model = spacy.load("en_core_web_sm",
                                 exclude=["tagger", "parser", "ner", "attribute_ruler", "lemmatizer"])
    
        keyBERT_model = keybert.KeyBERT(model= spacy_model)
    
        keywords = keyBERT_model.extract_keywords(preprocessed_post,
                                                  keyphrase_ngram_range=(1, 3),
                                                  stop_words="english")
        return keywords

    print(run_keyBERT_With_spacy())
    return (keybert,)


@app.cell
def _(keybert, preprocessed_post):
    def run_keyBERT_with_sentence_transformer():
        keyBERT_model = keybert.KeyBERT(model="all-MiniLM-L6-v2")

        keywords = keyBERT_model.extract_keywords(preprocessed_post,
                                                  keyphrase_ngram_range=(1, 3),
                                                  top_n=5,
                                                  stop_words="english")
        return keywords

    print(run_keyBERT_with_sentence_transformer())
    return


if __name__ == "__main__":
    app.run()
